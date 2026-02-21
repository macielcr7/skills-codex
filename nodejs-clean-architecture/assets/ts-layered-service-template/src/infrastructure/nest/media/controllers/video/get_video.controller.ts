import { BadRequestException, Controller, Get, Inject, NotFoundException, Param } from '@nestjs/common'

import { GetVideo } from '#/application/usecase/media/video/get_video/get_video'
import { VideoNotFoundError } from '#/domain/repository/media/video/video_not_found.error'

import { GetVideoParamsSchema } from '#/infrastructure/nest/media/controllers/video/get_video.schema'

@Controller('/media/videos')
export class GetVideoController {
  constructor(@Inject(GetVideo) private readonly getVideo: GetVideo) {}

  @Get(':id')
  async handle(@Param('id') id: string): Promise<{ id: string; title: string }> {
    const parsedParams = GetVideoParamsSchema.safeParse({ id })
    if (!parsedParams.success) {
      throw new BadRequestException(parsedParams.error.flatten())
    }

    try {
      return await this.getVideo.execute({ id: parsedParams.data.id })
    } catch (error) {
      if (error instanceof VideoNotFoundError) {
        throw new NotFoundException({ message: error.message })
      }
      throw error
    }
  }
}
