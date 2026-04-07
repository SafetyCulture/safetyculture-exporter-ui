export function trim(value: string, size: number): string {
    if (value.length > size) {
        return value.substring(0, size).concat(" ...")
    }
    return value
}

export function isNullOrEmptyObject(obj: Record<string, unknown> | null): boolean {
    if(obj === null) {
        return true
    }
    return Object.keys(obj).length === 0;
}

export const allTables: string[] = ['inspections', 'inspection_items', 'schedules', 'templates', 'template_permissions',
    'sites', 'site_members', 'groups', 'group_users', 'schedule_assignees', 'schedule_occurrences', 'actions',
    'action_assignees', 'action_timeline_items', 'issues', 'issue_timeline_items', 'assets', 'users', 'issue_assignees'];
