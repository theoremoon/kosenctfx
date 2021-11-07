import useSWR from "swr";

export interface Attachment {
  name: string;
  url: string;
}

export interface SolvedBy {
  solved_at: number;
  team_id: number;
  team_name: string;
}

export interface Task {
  id: number;
  name: string;
  description: string;
  flag: string;
  author: string;
  category: string;
  score: number;
  tags: string[];
  attachments: Attachment[];
  solved_by: SolvedBy[];

  is_open: boolean;
  is_survey: boolean;
}

const useAdminTasks = () => useSWR<Task[]>("/admin/list-challenges");

export default useAdminTasks;
