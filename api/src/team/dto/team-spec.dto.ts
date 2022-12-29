export class TeamSpecDto {
  teamName: string
  configRepo: TeamConfigRepoDto
}

export class TeamConfigRepoDto {
  source: string
  path: string
}
