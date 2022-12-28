export class SpecDto {
  teamName: string
  envName: string
  components: SpecComponentDto[]
}

export class SpecComponentDto {
  name: string
  type: string
  dependsOn: string[]
}
