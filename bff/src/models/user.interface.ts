export interface UserInterface{
  id: string;
  name: string;
  username?: string;
  groups?: string[];
}

export interface User {
  id: number;
  username: string;
  email: string;
  termAgreementStatus: boolean;
  role: string;
  archived: boolean;
  organizations: Organization[]
  created: Date;
  updated: Date;
}

export interface Organization {
  id: number
  name: string
}
