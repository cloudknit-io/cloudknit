import { ApiProperty } from "@nestjs/swagger";

export class DagDto {
  @ApiProperty()
  components: DagItemDto[];
}

export class DagItemDto {
  @ApiProperty()
  componentId: number;
  
  @ApiProperty()
  parentComponentIds: number[];
}
