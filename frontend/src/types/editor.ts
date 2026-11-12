export type EditorCommand =
  | 'h2'
  | 'h3'
  | 'bold'
  | 'italic'
  | 'strike'
  | 'quote'
  | 'ul'
  | 'ol'
  | 'link'
  | 'image'
  | 'code'
  | 'hr'

export interface EditorImagePayload {
  url: string
  alt: string
  caption: string
}
