import { IsNotEmpty, IsString } from "class-validator";
import { ComponentReconcile } from "src/typeorm";

export interface ComponentReconcileWrap extends ComponentReconcile {
    duration: number;
}

export class ApprovedByDto {
    @IsString()
    @IsNotEmpty()
    email: string;
}
