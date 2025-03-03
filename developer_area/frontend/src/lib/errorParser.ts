import { isAPIError } from './client'

export type EncoreValidationErrors = {
  [key: string]: string[]
}

export function parseEncoreError(message: string): EncoreValidationErrors {
  const errors: EncoreValidationErrors = {}

  if (!message.startsWith('validation failed:')) {
    return errors
  }

  const cleanMessage = message.replace('validation failed: ', '')
  const validationErrors = cleanMessage.split('\n').filter((msg) => msg.trim())

  validationErrors.forEach((errorMsg) => {
    const errorInfo = parseValidationErrorMessage(errorMsg)
    if (!errorInfo) return

    const { field, tag } = errorInfo

    if (!errors[field]) {
      errors[field] = []
    }

    if (!errors[field].includes(tag)) {
      errors[field].push(tag)
    }
  })

  return errors
}

function parseValidationErrorMessage(message: string):
  | {
      field: string
      tag: string
    }
  | undefined {
  const fieldMatch = message.match(/Key: 'RegistrationParams\.(\w+)'/)
  if (!fieldMatch) return

  const field = fieldMatch[1]

  const tagMatch = message.match(/'(\w+)' tag$/)
  if (!tagMatch) return

  const tag = tagMatch[1]

  return { field, tag }
}

export const mapEncoreErrorToFormErrors = (
  error: unknown,
  mappings: Record<string, Record<string, string>>,
): EncoreValidationErrors | undefined => {
  if (!isAPIError(error)) {
    console.error('Unexpected error:', error)
    return
  }

  const validationErrors = parseEncoreError(error.message)
  console.log('🪵 validationErrors', JSON.stringify(validationErrors))

  const userFriendlyErrors: EncoreValidationErrors = {}

  Object.entries(validationErrors).forEach(([field, tags]) => {
    userFriendlyErrors[field] = tags
      .map((tag) => mappings[field]?.[tag])
      .filter((msg): msg is string => msg !== undefined)
  })
  return userFriendlyErrors
}
