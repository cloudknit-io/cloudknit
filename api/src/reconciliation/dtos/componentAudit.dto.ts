import { IsNotEmpty, IsString } from "class-validator";

export interface ComponentAudit {
    reconcileId: number;
    duration: number;
    status: string;
    startDateTime: string;
    approvedBy?: string;
}

export class ApprovedByDto {
    @IsString()
    @IsNotEmpty()
    email: string;
}
