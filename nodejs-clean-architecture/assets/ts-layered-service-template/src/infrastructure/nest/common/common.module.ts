import { Module } from '@nestjs/common'

import { HealthController } from '#/infrastructure/nest/common/controllers/health.controller'
import { UUIDGenerator } from '#/infrastructure/nest/common/services/uuid_generator'

@Module({
  controllers: [HealthController],
  providers: [UUIDGenerator],
  exports: [UUIDGenerator],
})
export class CommonModule {}

