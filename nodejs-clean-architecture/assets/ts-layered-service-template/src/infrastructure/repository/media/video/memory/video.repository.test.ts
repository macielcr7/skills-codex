import { describe, expect, it } from 'vitest'

import { Video } from '#/domain/entity/media/video/video'
import { MemoryVideoRepository } from '#/infrastructure/repository/media/video/memory/video.repository'

describe('MemoryVideoRepository', () => {
  it('saves and returns by id', async () => {
    const repo = new MemoryVideoRepository()
    const video = new Video({
      id: '550e8400-e29b-41d4-a716-446655440000',
      title: 'hello',
    })
    video.validate()

    await repo.save(video)

    const found = await repo.getByID(video.id)
    expect(found).toBe(video)
  })

  it('returns null when missing', async () => {
    const repo = new MemoryVideoRepository()
    const found = await repo.getByID('550e8400-e29b-41d4-a716-446655440000')
    expect(found).toBeNull()
  })
})

