export class EnvSpecDto {
  envName: string
  components: EnvSpecComponentDto[]
}

export class EnvSpecComponentDto {
  name: string
  type: string
  dependsOn: string[]
}
