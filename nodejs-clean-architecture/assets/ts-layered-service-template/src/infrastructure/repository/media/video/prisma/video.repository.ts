import type { PrismaClient } from '@prisma/client'

import { Video } from '#/domain/entity/media/video/video'
import type { VideoRepository } from '#/domain/repository/media/video/video.repository.interface'

export class PrismaVideoRepository implements VideoRepository {
  constructor(private readonly prisma: PrismaClientLike) {}

  async save(video: Video): Promise<void> {
    await this.prisma.video.upsert({
      where: { id: video.id },
      update: { title: video.title },
      create: { id: video.id, title: video.title },
    })
  }

  async getByID(id: string): Promise<Video | null> {
    const row = await this.prisma.video.findUnique({ where: { id } })
    if (!row) {
      return null
    }

    const video = new Video({ id: row.id, title: row.title })
    video.validate()
    return video
  }
}

export type PrismaClientLike = Pick<PrismaClient, 'video'>
