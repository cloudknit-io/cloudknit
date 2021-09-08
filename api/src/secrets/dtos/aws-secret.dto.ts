
interface AwsSecret {
    pathName: string;
    value: string;
}

export interface AwsSecretDto {
    secrets: AwsSecret[];
}