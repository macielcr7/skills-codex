import { describe, expect, it, vi } from 'vitest'

import { Video } from '#/domain/entity/media/video/video'
import { PrismaVideoRepository } from '#/infrastructure/repository/media/video/prisma/video.repository'

describe('PrismaVideoRepository', () => {
  it('save() calls prisma upsert', async () => {
    const prisma = {
      video: {
        upsert: vi.fn().mockResolvedValue(undefined),
        findUnique: vi.fn(),
      },
    }

    const repo = new PrismaVideoRepository(prisma as never)
    const video = new Video({
      id: '550e8400-e29b-41d4-a716-446655440000',
      title: 'hello',
    })
    video.validate()

    await repo.save(video)

    expect(prisma.video.upsert).toHaveBeenCalledWith({
      where: { id: video.id },
      update: { title: video.title },
      create: { id: video.id, title: video.title },
    })
  })

  it('getByID() returns null when missing', async () => {
    const prisma = {
      video: {
        upsert: vi.fn(),
        findUnique: vi.fn().mockResolvedValue(null),
      },
    }

    const repo = new PrismaVideoRepository(prisma as never)
    const found = await repo.getByID('550e8400-e29b-41d4-a716-446655440000')
    expect(found).toBeNull()
  })

  it('getByID() returns validated Video when found', async () => {
    const prisma = {
      video: {
        upsert: vi.fn(),
        findUnique: vi.fn().mockResolvedValue({
          id: '550e8400-e29b-41d4-a716-446655440000',
          title: 'hello',
        }),
      },
    }

    const repo = new PrismaVideoRepository(prisma as never)
    const found = await repo.getByID('550e8400-e29b-41d4-a716-446655440000')
    expect(found?.id).toBe('550e8400-e29b-41d4-a716-446655440000')
    expect(found?.title).toBe('hello')
  })
})

