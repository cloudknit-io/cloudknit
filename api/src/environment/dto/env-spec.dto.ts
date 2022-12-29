export class EnvSpecDto {
  teamName: string
  envName: string
  components: EnvSpecComponentDto[]
}

export class EnvSpecComponentDto {
  name: string
  type: string
  dependsOn: string[]
}
