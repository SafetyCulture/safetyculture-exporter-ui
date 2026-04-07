import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export type WithElementRef<T, El extends HTMLElement = HTMLElement> = T & { ref?: El | null };
export type WithoutChild<T> = T extends { child?: unknown } ? Omit<T, 'child'> : T;
export type WithoutChildrenOrChild<T> = T extends { children?: unknown; child?: unknown }
  ? Omit<T, 'children' | 'child'>
  : T extends { children?: unknown }
    ? Omit<T, 'children'>
    : T extends { child?: unknown }
      ? Omit<T, 'child'>
      : T;

export function trim(value: string, size: number): string {
  if (value.length > size) {
    return value.substring(0, size).concat(' ...');
  }
  return value;
}

export function isNullOrEmptyObject(obj: Record<string, unknown> | null): boolean {
  if (obj === null) {
    return true;
  }
  return Object.keys(obj).length === 0;
}

export const allTables: string[] = [
  'inspections',
  'inspection_items',
  'schedules',
  'templates',
  'template_permissions',
  'sites',
  'site_members',
  'groups',
  'group_users',
  'schedule_assignees',
  'schedule_occurrences',
  'actions',
  'action_assignees',
  'action_timeline_items',
  'issues',
  'issue_timeline_items',
  'assets',
  'users',
  'issue_assignees',
];
