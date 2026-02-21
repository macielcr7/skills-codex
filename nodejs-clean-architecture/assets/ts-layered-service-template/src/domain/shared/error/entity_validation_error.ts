import type { ZodIssue } from 'zod'

export class EntityValidationError extends Error {
  public readonly issues: ZodIssue[]

  constructor(entity: string, issues: ZodIssue[]) {
    super(`${entity} validation failed`)
    this.name = 'EntityValidationError'
    this.issues = issues
  }
}

