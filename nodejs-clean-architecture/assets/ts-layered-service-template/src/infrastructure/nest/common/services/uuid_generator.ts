import { randomUUID } from 'node:crypto'

import { Injectable } from '@nestjs/common'

import type { IdGenerator } from '#/application/service/shared/id_generator.interface'

@Injectable()
export class UUIDGenerator implements IdGenerator {
  newID(): string {
    return randomUUID()
  }
}
