import Loading from "components/loading";
import useMessage from "lib/useMessage";
import React, { useCallback, useState } from "react";
import useAdminTasks, { Task } from "lib/api/admin/tasks";
import { useRouter } from "next/router";
import AdminLayout from "components/adminLayout";
import { api } from "lib/api";

type taskElementProps = {
  task: Task;
  isOpened: boolean;
  onUpdateOpened: (opened: boolean) => void;
  onClickCallback: () => void;
};

const TaskElement = ({
  task,
  isOpened,
  onUpdateOpened,
  onClickCallback,
}: taskElementProps) => {
  return (
    <>
      <td title={task.name}>{task.name}</td>
      <td>{task.score}</td>
      <td>{task.solved_by.length}</td>
      <td>{task.category}</td>
      <td>{task.author}</td>
      <td>{task.deployment}</td>
      <td className="form-switch">
        <input
          className="form-check-input"
          type="checkbox"
          role="switch"
          style={{ margin: 0 }}
          checked={isOpened}
          onChange={() => {
            onUpdateOpened(!isOpened);
          }}
        />
      </td>
      <td>{task.is_survey && "🗒️"}</td>
      <td style={{ cursor: "pointer" }} onClick={onClickCallback}>
        👀
      </td>
      <td>
        <code title={task.flag}>
          <pre>{task.flag}</pre>
        </code>
      </td>
    </>
  );
};

interface TasksProps {
  tasks: Task[];
}

interface TasksMD {
  text: string;
}

const Tasks = ({ tasks }: TasksProps) => {
  const { mutate } = useAdminTasks();
  const { message, error, text: textMessage } = useMessage();

  const openState = new Map(tasks.map((t) => [t.id, t.is_open]));
  const [taskOpenState, setTaskOpenState] = useState(openState);
  const router = useRouter();

  const openCloseChallenge = useCallback(() => {
    tasks.forEach(async (t) => {
      const next_state = taskOpenState.get(t.id);
      if (t.is_open === next_state) {
        return;
      }

      const endpoint = next_state
        ? "/admin/open-challenge"
        : "/admin/close-challenge";
      try {
        const res = await api.post(endpoint, {
          name: t.name,
        });
        message(res);
      } catch (e) {
        error(e);
      }
    });
    mutate();
  }, [api, message, error, tasks]);

  const generateTasksMD = useCallback(async () => {
    const res = await api.get<TasksMD>("/admin/tasks.md");
    navigator.clipboard.writeText(res.data.text);
    textMessage("Copied tasks.md to clipboard");
  }, [api, textMessage]);

  const generateCheckPortPy = useCallback(async () => {
    const res = await api.get<TasksMD>("/admin/check-port.py");

    const link = document.createElement("a");
    link.href =
      "data:text/plain;charset=utf-8," + encodeURIComponent(res.data.text);
    link.download = "check-port.py";
    document.body.appendChild(link);
    link.click();
    setTimeout(() => {
      if (link.parentNode) {
        link.parentNode.removeChild(link);
      }
    }, 1000);
  }, [api]);

  return (
    <>
      <h5 className="mt-4">Download Scripts</h5>
      <div>
        <button
          type="button"
          className="btn btn-primary mx-2"
          onClick={generateTasksMD}
        >
          tasks.md
        </button>
        <button
          type="button"
          className="btn btn-primary"
          onClick={generateCheckPortPy}
        >
          check-port.py
        </button>
      </div>

      <h5 className="mt-4">Tasks</h5>
      <table className="table table-responsive">
        <thead>
          <tr>
            <th>Name</th>
            <th>Score</th>
            <th>#Solve</th>
            <th>Category</th>
            <th>Author</th>
            <th>Deployment Type</th>
            <th>Is Open?</th>
            <th>Is Survey?</th>
            <th>Preview</th>
            <th>Flag</th>
          </tr>
        </thead>
        <tbody>
          {tasks.map((task) => (
            <tr key={task.name}>
              <TaskElement
                task={task}
                isOpened={taskOpenState.get(task.id) || false}
                onUpdateOpened={(isOpen) => {
                  setTaskOpenState((prev) =>
                    new Map(prev).set(task.id, isOpen)
                  );
                }}
                onClickCallback={() => {
                  router.push(`/admin/tasks/${task.id}`, undefined, {
                    shallow: true,
                    scroll: false,
                  });
                }}
              />
            </tr>
          ))}
        </tbody>
      </table>

      <button
        type="button"
        className="btn btn-primary"
        onClick={openCloseChallenge}
      >
        Open/Close Challenges
      </button>
    </>
  );
};

const TasksDefault = () => {
  const { data: tasks } = useAdminTasks();

  if (!tasks) {
    return <Loading />;
  }
  return <Tasks tasks={tasks} />;
};

TasksDefault.getLayout = AdminLayout;

export default TasksDefault;
