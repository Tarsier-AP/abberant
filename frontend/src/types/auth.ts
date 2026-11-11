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
