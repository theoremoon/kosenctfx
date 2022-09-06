import { useRouter } from "next/router";
import Loading from "../../../components/loading";
import parentpath from "../../../lib/parentpath";
import useAdminTasks from "../../../lib/api/admin/tasks";
import TaskModalView from "theme/taskmodal";
import { FlagSubmitParams } from "props/taskmodal";
import { useForm } from "react-hook-form";
import useConfig from "lib/api/admin/config";

const Task = () => {
  const router = useRouter();
  const { id } = router.query;
  const { register } = useForm<FlagSubmitParams>();

  const { data: tasks } = useAdminTasks();

  // admin-check
  const { data: config, error: configError } = useConfig();
  if (config === undefined || configError !== undefined) {
    return <></>;
  }

  if (!tasks) {
    return <Loading />;
  }

  const filterdTasks = tasks.filter((t) => t.id === Number(id));
  if (filterdTasks.length !== 1) {
    return <Loading />;
  }
  const task = filterdTasks[0];
  const tasksPath = parentpath(router.pathname); // /admin/tasks

  return (
    <>
      <TaskModalView
        task={task}
        onClose={() =>
          router.push(tasksPath, undefined, {
            scroll: false,
            shallow: true,
          })
        }
        registerFlag={register}
        onFlagSubmit={() => {
          // noop
        }}
        isSolved={false}
      />
    </>
  );
};

export default Task;
