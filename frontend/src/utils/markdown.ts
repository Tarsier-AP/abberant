import DOMPurify from 'dompurify'
import { marked, type Tokens } from 'marked'

import { buildSrcset } from '@/api/media'

function escapeAttribute(value: string): string {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

/** Responsive hints for body images: full width on phones, capped on desktop. */
const BODY_SIZES = '(max-width: 767px) 100vw, 828px'

function figureFromImage(token: Tokens.Image): string {
  const href = escapeAttribute(token.href ?? '')
  const alt = escapeAttribute(token.text ?? '')
  const caption = token.title ? escapeAttribute(token.title) : ''
  const srcset = buildSrcset(token.href ?? '')
  const srcsetAttr = srcset ? ` srcset="${escapeAttribute(srcset)}" sizes="${BODY_SIZES}"` : ''
  // Vertical images are shown whole over a blurred backdrop of the same
  // photo (`article.css` + `detectPortraitFigures`); the fixed 16:9 frame
  // keeps the layout stable before anything loads. The backdrop reuses the
  // same URL, so the browser downloads the image only once.
  const bg = `<img class="figure__bg" src="${href}" alt="" aria-hidden="true" decoding="async">`
  const img = `<img class="figure__img" src="${href}"${srcsetAttr} alt="${alt}" width="1200" height="675" decoding="async">`
  const media = `<span class="figure__media">${bg}${img}</span>`

  return caption
    ? `<figure>${media}<figcaption>${caption}</figcaption></figure>`
    : `<figure>${media}</figure>`
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
