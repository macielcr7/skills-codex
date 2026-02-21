import { describe, expect, it } from 'vitest'

import { GetVideoParamsSchema } from '#/infrastructure/nest/media/controllers/video/get_video.schema'

describe('GetVideoParamsSchema', () => {
  it('parses valid params', () => {
    const result = GetVideoParamsSchema.safeParse({
      id: '550e8400-e29b-41d4-a716-446655440000',
    })
    expect(result.success).toBe(true)
  })

  it('rejects invalid params', () => {
    const result = GetVideoParamsSchema.safeParse({ id: 'invalid' })
    expect(result.success).toBe(false)
  })
})
