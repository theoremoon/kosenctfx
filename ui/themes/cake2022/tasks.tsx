import TaskCard from "./components/taskCard";
import { TasksProps } from "props/tasks";

const Tasks = ({ tasks, sortTask, isSolved }: TasksProps) => {
  return (
    <div
      style={{
        display: "grid",
        gap: "1rem",
        gridTemplateColumns: "repeat(auto-fit, 12rem)",
      }}
    >
      {sortTask(tasks).map((t) => {
        const solved = isSolved(t);
        return <TaskCard key={t.id} task={t} isSolved={solved} />;
      })}
    </div>
  );
};

export default Tasks;
