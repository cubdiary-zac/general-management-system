export type TaskStatus = 'todo' | 'in_progress' | 'in_review' | 'done'

export const orderedStatuses: TaskStatus[] = ['todo', 'in_progress', 'in_review', 'done']

export function nextTaskStatus(status: TaskStatus): TaskStatus | null {
  const idx = orderedStatuses.indexOf(status)
  if (idx < 0 || idx === orderedStatuses.length - 1) {
    return null
  }
  return orderedStatuses[idx + 1]
}
