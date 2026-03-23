const GQL_ENDPOINT = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/query'

export async function gql<T>(
  query: string,
  variables?: Record<string, unknown>,
  sessionId?: string | null,
): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (sessionId) headers['Authorization'] = `Bearer ${sessionId}`

  const res = await fetch(GQL_ENDPOINT, {
    method: 'POST',
    headers,
    body: JSON.stringify({ query, variables }),
  })

  const json = await res.json()
  if (json.errors?.length) throw new Error(json.errors[0].message)
  return json.data as T
}
