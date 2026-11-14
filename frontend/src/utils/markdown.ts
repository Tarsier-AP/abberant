import DOMPurify from 'dompurify'
import { marked, type Tokens } from 'marked'

function escapeAttribute(value: string): string {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function figureFromImage(token: Tokens.Image): string {
  const href = escapeAttribute(token.href ?? '')
  const alt = escapeAttribute(token.text ?? '')
  const caption = token.title ? escapeAttribute(token.title) : ''
  // Fixed 16:9 intrinsic dimensions keep layout stable before the image
  // loads; the frame is enforced by CSS (`article.css`) too.
  const img = `<img src="${href}" alt="${alt}" width="1200" height="675" decoding="async">`

  return caption
    ? `<figure>${img}<figcaption>${caption}</figcaption></figure>`
    : `<figure>${img}</figure>`
}

/**
 * Images with a markdown title (`![alt](src "caption")`) render as a figure
 * with a visible caption; the title is never used as a tooltip. A paragraph
 * holding only an image renders the figure directly, without the wrapping
 * paragraph.
 */
marked.use({
  gfm: true,
  renderer: {
    image(token) {
      return figureFromImage(token)
    },
    paragraph(token) {
      const tokens = token.tokens ?? []
      const first = tokens[0]

      if (tokens.length === 1 && first && first.type === 'image') {
        return figureFromImage(first as Tokens.Image)
      }

      // Fall back to the default paragraph rendering.
      return false
    },
  },
})

/**
 * Renders markdown source to sanitized HTML. Every consumer (editor preview,
 * article page) must go through this function — the output is safe for
 * `v-html`.
 */
export function renderMarkdown(source: string): string {
  const html = marked.parse(source, { async: false })

  return DOMPurify.sanitize(html)
}
