export type ConfigRecord = ClientConfig | UserConfig;

export interface ConfigData{
    [key: string]: ConfigRecord;
}

export interface ClientConfig {
    client_id:     string;
    client_secret: string | undefined;
    account:        AccountInfo;
}

export interface UserConfig {
    subject: string | undefined;
    account: AccountInfo;
}

export interface AccountInfo {
    id: number;
    group: string;
    role:  string;
}
