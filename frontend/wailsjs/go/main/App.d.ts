// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {api} from '../models';
import {main} from '../models';

export function CancelExport():Promise<void>;

export function CheckDBConnection():Promise<void>;

export function ExportCSV():Promise<void>;

export function ExportReports():Promise<void>;

export function ExportSQL():Promise<void>;

export function ExportSQLite():Promise<void>;

export function GetSettingDir():Promise<string>;

export function GetSettings():Promise<api.ExporterConfiguration>;

export function GetTemplates():Promise<Array<api.TemplateResponseItem>>;

export function GetUserHomeDirectory():Promise<string>;

export function OpenDirectory(arg1:string):Promise<void>;

export function ReadBuild():Promise<string>;

export function ReadExportStatus():Promise<void>;

export function ReadVersion():Promise<main.VersionResponse>;

export function SaveSettings(arg1:api.ExporterConfiguration):Promise<void>;

export function SelectDirectory(arg1:string):Promise<string>;

export function TriggerUpdate(arg1:string):Promise<boolean>;

export function ValidateApiKey(arg1:string):Promise<string>;

export function ValidateExportDirectory():Promise<boolean>;
