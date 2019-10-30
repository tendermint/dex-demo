export type ActionType<Payload> = {
  type: string
  payload: Payload
  error?: boolean
  meta?: any
}
