import { describe, expect, it } from 'vitest'

import { UploadVideoBodySchema } from '#/infrastructure/nest/media/controllers/video/upload_video.schema'

describe('UploadVideoBodySchema', () => {
  it('parses valid body', () => {
    const result = UploadVideoBodySchema.safeParse({ title: 'hello' })
    expect(result.success).toBe(true)
  })

  it('rejects invalid body', () => {
    const result = UploadVideoBodySchema.safeParse({ title: '' })
    expect(result.success).toBe(false)
  })
})
