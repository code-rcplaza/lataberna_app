<script setup lang="ts">
import { useAuthStore } from '@/stores/useAuthStore'

const auth = useAuthStore()
</script>

<template>
  <div class="max-w-3xl mx-auto px-8 py-20 flex flex-col gap-16">

    <!-- Hero -->
    <section class="flex flex-col gap-6">
      <p class="font-label uppercase tracking-widest text-xs font-semibold text-secondary">
        D&amp;D 5e · Generador de personajes
      </p>
      <h1 class="text-5xl font-bold leading-tight text-on-surface">
        Tu próximo personaje,<br>
        <span class="text-primary italic">en segundos.</span>
      </h1>
      <p class="text-on-surface-variant text-lg leading-relaxed max-w-xl">
        Nombre, estadísticas, trasfondo y motivación — coherentes entre sí, en español nativo.
        Listo para sentarse a jugar.
      </p>
      <div class="flex flex-wrap gap-4 pt-2">
        <RouterLink
          to="/forja"
          class="inline-flex items-center gap-2 px-6 py-3 bg-primary text-on-primary font-label font-bold uppercase tracking-widest text-xs hover:bg-primary-container hover:text-on-primary-container transition-colors"
        >
          <span class="material-symbols-outlined text-sm">auto_fix_high</span>
          Generar personaje
        </RouterLink>
        <RouterLink
          v-if="auth.isAuthenticated"
          to="/biblioteca"
          class="inline-flex items-center gap-2 px-6 py-3 border border-outline-variant text-on-surface-variant font-label font-bold uppercase tracking-widest text-xs hover:border-outline hover:text-on-surface transition-colors"
        >
          <span class="material-symbols-outlined text-sm">library_books</span>
          Mi biblioteca
        </RouterLink>
        <RouterLink
          v-else
          to="/auth"
          class="inline-flex items-center gap-2 px-6 py-3 border border-outline-variant text-on-surface-variant font-label font-bold uppercase tracking-widest text-xs hover:border-outline hover:text-on-surface transition-colors"
        >
          <span class="material-symbols-outlined text-sm">lock_open</span>
          Iniciar sesión
        </RouterLink>
      </div>
    </section>

    <!-- Divisor -->
    <div class="border-t border-outline-variant/30" />

    <!-- Features -->
    <section class="grid grid-cols-1 sm:grid-cols-3 gap-8">
      <div class="flex flex-col gap-3">
        <span class="material-symbols-outlined text-primary text-3xl">casino</span>
        <h3 class="font-headline text-on-surface text-lg">Todo aleatorio</h3>
        <p class="text-on-surface-variant text-sm leading-relaxed">
          Sin configurar nada obtenés un personaje completo. Un click, listo para jugar.
        </p>
      </div>
      <div class="flex flex-col gap-3">
        <span class="material-symbols-outlined text-primary text-3xl">lock</span>
        <h3 class="font-headline text-on-surface text-lg">Control total</h3>
        <p class="text-on-surface-variant text-sm leading-relaxed">
          Bloqueá lo que te gusta y regenerá el resto. Nombre, stats y narrativa de forma independiente.
        </p>
      </div>
      <div class="flex flex-col gap-3">
        <span class="material-symbols-outlined text-primary text-3xl">library_books</span>
        <h3 class="font-headline text-on-surface text-lg">Tu biblioteca</h3>
        <p class="text-on-surface-variant text-sm leading-relaxed">
          Guardá tus personajes y NPCs favoritos. Accedé, editá y reutilizalos en cualquier sesión.
        </p>
      </div>
    </section>

  </div>
</template>
