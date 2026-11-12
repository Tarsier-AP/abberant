export interface LoginCredentials {
  login: string
  password: string
}

export interface RegisterCredentials {
  login: string
  username: string
  acceptedTerms: boolean
  acceptedPolicy: boolean
}

export type TwoFactorMethod = 'authenticator' | 'sms' | 'email'

export type UserRole = 'reader' | 'editor' | 'owner'

export interface AuthAccount {
  login: string
  username: string
  avatar: string
  additionalEmail: string
  phone: string
  activated: boolean
  twoFactorEnabled: boolean
  twoFactorMethod: TwoFactorMethod | ''
  smsNotifications: boolean
  emailNotifications: boolean
  role: UserRole
}
