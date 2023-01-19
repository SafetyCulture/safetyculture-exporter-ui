export function trim(value, size) {
    if (value.length > size) {
        return value.substring(0, size).concat(" ...")
    }
    return value
}
export const allTables = ['inspections', 'inspection_items', 'schedules', 'templates', 'template_permissions',
    'sites', 'site_members', 'groups', 'group_users', 'schedule_assignees', 'schedule_occurrences', 'actions',
    'action_assignees', 'issues', 'assets', 'users'];
