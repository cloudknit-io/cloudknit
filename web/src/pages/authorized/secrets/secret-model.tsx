export interface SecretModel {
    name: string;
    key: string;
    notRequired?: boolean;
    multiline?: boolean;
    hide?: boolean;
    immutable?: boolean;
}