<script lang="ts">
    import { shadowConfig } from "../lib/store";
    import { push } from "@keenmate/svelte-spa-router";
    import Button from "../components/Button.svelte";
    import StatusBar from "../components/StatusBar.svelte";

    interface DataItem {
        id: string;
        name: string;
    }

    interface DataRow {
        left: DataItem | null;
        right: DataItem | null;
    }

    let data: DataRow[] = [
        {
            "left": { "id": "inspections", "name": "Inspections" },
            "right": { "id": "inspection_items", "name": "Inspection Items"},
        },
        {
            "left":  { "id": "templates", "name": "Templates" },
            "right": { "id": "template_permissions", "name": "Template Permissions"}
        },
        {
            "left":  { "id": "sites", "name": "Sites"},
            "right": { "id": "site_members", "name": "Site Members"}
        },
        {
            "left":  { "id": "groups", "name": "Groups"},
            "right": { "id": "group_users", "name": "Group Users"}
        },
        {
            "left":  { "id": "users", "name": "Users"},
            "right": { "id": "schedules", "name": "Schedules"}
        },
        {
            "left":  { "id": "schedule_assignees", "name": "Schedule Assignees"},
            "right": { "id": "schedule_occurrences", "name": "Schedule Occurrences"}
        },
        {
            "left":  { "id": "actions", "name": "Actions"},
            "right": { "id": "issues", "name": "Issues"}
        },
        {
            "left":  { "id": "action_timeline_items", "name": "Action Timeline Items"},
            "right":  { "id": "issue_timeline_items", "name": "Issue Timeline Items"},
        },
        {
            "left":  { "id": "action_assignees", "name": "Action Assignees"},
            "right": { "id": "issue_assignees", "name": "Issue Assignees"}
        },
        {
            "left": { "id": "assets", "name": "Assets"},
            "right": { "id": "training_course_progresses", "name": "Training course completions"},
        }
    ]

    function trimName(org: string): string {
        if (org.length > 80) {
            return org.substring(0, 80).concat(" ...")
        }
        return org
    }

    let isChecked = $state(false);
    if (($shadowConfig as any)["Export"]["Tables"].length === 0) {
        let all: string[] = []
        data.forEach(function (e) {
            if (e.left !== null) {
                all.push(e.left.id)
            }
            if (e.right !== null) {
                all.push(e.right.id)
            }
        });
        ($shadowConfig as any)["Export"]["Tables"] = all
        isChecked = true;
    }

    function toggleHeaderCheckbox() {
        if (isChecked) {
            isChecked = false;
        }
    }

    function toggleBodyCheckboxes() {
        const checkboxes = document.querySelectorAll('.table-body input[type="checkbox"]') as NodeListOf<HTMLInputElement>;
        for (const checkbox of checkboxes) {
            checkbox.checked = !isChecked;
        }
    }

    function handleDone() {
        let selectedTables: string[] = [];

        const checkboxes = document.querySelectorAll('.table-body input[type="checkbox"]') as NodeListOf<HTMLInputElement>;
        for (const checkbox of checkboxes) {
            if (checkbox.checked) {
                selectedTables.push((checkbox as any).__value)
            }
        }

        let maxData = 0
        data.forEach(function (e) {
            if (e.left !== null) maxData++
            if (e.right !== null) maxData++
        });

        if (selectedTables.length === maxData) {
            ($shadowConfig as any)["Export"]["Tables"] = []
        } else {
            ($shadowConfig as any)["Export"]["Tables"] = selectedTables
        }

        push("/config")
    }
</script>

<div class="px-8 pt-8">
    <section class="flex items-center justify-between">
        <div class="flex items-center">
            <div class="my-4 font-semibold text-2xl">Data set selection</div>
        </div>
        <div class="flex items-center">
            <Button label="Done" type="active-white" onClick={handleDone}/>
        </div>
    </section>

    <section class="mt-4">
        <div class="bg-[#DBDFEB] text-text-gray">
            <div class="flex h-9 items-center px-2">
                <input type="checkbox" class="size-5 accent-primary" onclick={toggleBodyCheckboxes} bind:checked={isChecked}/>
                <div class="ml-8">Data set table</div>
            </div>
        </div>
        <div class="table-body mt-2 overflow-y-hidden text-text-gray" style="-ms-overflow-style: none; scrollbar-width: none;">
        {#each data as { left, right }}
            <div class="flex h-13 w-full">
                {#if left}
                <div class="flex w-1/2 items-center">
                    <input type="checkbox" class="size-5 accent-primary" onclick={toggleHeaderCheckbox} bind:group={($shadowConfig as any)["Export"]["Tables"]} value={left.id}/>
                    <img class="ml-8" src="../images/template-icon.svg" alt="template"/>
                    <div class="ml-2">{trimName(left.name)}</div>
                </div>
                {/if}

                {#if right}
                <div class="flex w-1/2 items-center">
                    <input type="checkbox" class="size-5 accent-primary" onclick={toggleHeaderCheckbox} bind:group={($shadowConfig as any)["Export"]["Tables"]} value={right.id}/>
                    <img class="ml-8" src="../images/template-icon.svg" alt="template"/>
                    <div class="ml-2">{trimName(right.name)}</div>
                </div>
                {/if}
            </div>
        {/each}
        </div>
    </section>
</div>

<StatusBar/>
