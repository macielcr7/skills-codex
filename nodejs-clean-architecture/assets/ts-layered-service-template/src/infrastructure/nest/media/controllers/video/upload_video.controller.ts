import { BadRequestException, Body, Controller, Inject, Post } from '@nestjs/common'

import { UploadVideo } from '#/application/usecase/media/video/upload_video/upload_video'
import { EntityValidationError } from '#/domain/shared/error/entity_validation_error'

import { UploadVideoBodySchema } from '#/infrastructure/nest/media/controllers/video/upload_video.schema'

@Controller('/media/videos')
export class UploadVideoController {
  constructor(
    @Inject(UploadVideo) private readonly uploadVideo: UploadVideo,
  ) {}

  @Post()
  async handle(@Body() body: unknown): Promise<{ id: string }> {
    const parsedBody = UploadVideoBodySchema.safeParse(body)
    if (!parsedBody.success) {
      throw new BadRequestException(parsedBody.error.flatten())
    }

    try {
      return await this.uploadVideo.execute({ title: parsedBody.data.title })
    } catch (error) {
      if (error instanceof EntityValidationError) {
        throw new BadRequestException({ message: error.message, issues: error.issues })
      }
      throw error
    }
  }
}
