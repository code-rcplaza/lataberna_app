package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"

	"forge-rpg/graph"
	"forge-rpg/graph/generated"
	infradb "forge-rpg/internal/infrastructure/db"
	"forge-rpg/internal/infrastructure/email"
	graphqlmw "forge-rpg/internal/infrastructure/graphql"
	"forge-rpg/internal/domain/ports"
	"forge-rpg/internal/usecase/auth"
	"forge-rpg/internal/usecase/character"
	"forge-rpg/internal/usecase/namegen"
	"forge-rpg/internal/usecase/narrativegen"
)

func main() {
	// Carga .env si existe — en producción las vars vienen del entorno directamente
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using environment variables")
	}

	dbPath := getEnv("DB_PATH", "forge.db")
	linkBase := getEnv("LINK_BASE", "http://localhost:8080/auth/verify")

	database, err := infradb.Open(dbPath)
	if err != nil {
		log.Fatalf("db.Open: %v", err)
	}
	defer database.Close()

	// Seed content on startup (idempotent — skips if already populated)
	ctx := context.Background()
	if err := infradb.SeedContentIfEmpty(ctx, database); err != nil {
		log.Fatalf("seed content: %v", err)
	}

	// Repositories
	userRepo := infradb.NewUserRepository(database)
	sessionRepo := infradb.NewSessionRepository(database)
	tokenRepo := infradb.NewTokenRepository(database)
	charRepo := infradb.NewCharacterRepository(database)
	narrativeRepo := infradb.NewNarrativeRepository(database)
	nameRepo := infradb.NewNameRepository(database)

	// Mailer
	var mailer ports.Mailer
	resendKey := getEnv("RESEND_API_KEY", "")
	if resendKey != "" {
		from := getEnv("RESEND_FROM", "La Taberna <noreply@lataberna.app>")
		mailer = email.NewResendMailer(resendKey, from)
		log.Println("Mailer: Resend (production)")
	} else {
		mailer = &logMailer{}
		log.Println("Mailer: stdout (development — set RESEND_API_KEY for real emails)")
	}

	// Services
	authSvc := auth.NewService(userRepo, sessionRepo, tokenRepo, mailer, linkBase)
	managerSvc := character.NewService(charRepo)
	narrativeSvc := narrativegen.New(narrativeRepo)
	nameSvc := namegen.New(nameRepo)
	creatorSvc := character.NewCreator(narrativeSvc, nameSvc)

	// GraphQL
	resolver := graph.NewResolver(authSvc, managerSvc, creatorSvc)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("FORGE RPG", "/query"))
	mux.Handle("/query", graphqlmw.AuthMiddleware(authSvc)(srv))

	port := getEnv("PORT", "8080")
	log.Printf("Server running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	allowed := getEnv("CORS_ORIGIN", "http://localhost:5173")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowed)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// logMailer prints magic links to stdout — replace with real SMTP in production.
type logMailer struct{}

func (m *logMailer) SendMagicLink(ctx context.Context, email, link string) error {
	log.Printf("[MAGIC LINK] To: %s | Link: %s", email, link)
	return nil
}
