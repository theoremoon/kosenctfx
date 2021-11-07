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
  category: string;
  author: string;
  score: number;
  tags: string[];
  attachments: Attachment[];
  solved_by: SolvedBy[];

  is_open: boolean;
  is_survey: boolean;
}

const useTasks = () => useSWR<Task[]>("/tasks");

export default useTasks;
