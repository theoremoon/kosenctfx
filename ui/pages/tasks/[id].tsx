import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import { isStaticMode, revalidateInterval } from "lib/static";
import { GetStaticPaths, GetStaticProps } from "next";
import { useRouter } from "next/router";
import Loading from "../../components/loading";
import useAccount, { Account, fetchAccount } from "../../lib/api/account";
import useTasks, { fetchTasks, Task } from "../../lib/api/tasks";
import parentpath from "../../lib/parentpath";
import useMessage from "lib/useMessage";
import { useForm, SubmitHandler } from "react-hook-form";
import { api } from "lib/api";
import { filterTask, sortTask, isSolved } from "lib/tasks";
import { useCallback } from "react";

import TasksView from "theme/tasks";
import TaskModalView from "theme/taskmodal";
import { FlagSubmitParams } from "props/taskmodal";

type taskProps = {
  taskID: number;
  tasks: Task[];
} & AllPageProps;

const TasksDefault = ({ taskID, tasks: defaultTasks }: taskProps) => {
  const router = useRouter();
  const { data: tasks, mutate } = useTasks(defaultTasks);
  const { data: account } = useAccount(null);
  const { message, error } = useMessage();
  const { register, handleSubmit } = useForm<FlagSubmitParams>();
  const isSolvedByTeam = useCallback(isSolved(account || null), [account]);
  if (!tasks || account === undefined) {
    return <Loading />;
  }

  const onSubmit: SubmitHandler<FlagSubmitParams> = async (values) => {
    try {
      const res = await api.post("/submit", {
        flag: values.flag,
      });
      message(res);
      mutate();
    } catch (e) {
      error(e);
    }
  };

  const filterdTasks = tasks.filter((t) => t.id === taskID);
  if (filterdTasks.length !== 1) {
    return <Loading />;
  }
  const task = filterdTasks[0];
  const tasksPath = parentpath(router.pathname); // /tasks

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
        onFlagSubmit={handleSubmit(onSubmit)}
        isSolved={isSolvedByTeam(task)}
      />
      <TasksView
        tasks={tasks}
        filterTask={filterTask}
        sortTask={sortTask}
        isSolved={isSolvedByTeam}
      />
    </>
  );
};

export const getStaticProps: GetStaticProps<taskProps> = async (context) => {
  const id = context.params?.id;
  const tasks = await fetchTasks();
  const ctf = await fetchCTF();
  return {
    props: {
      taskID: Number(id),
      tasks: tasks,
      ctf: ctf,
    },
    // staticModeでないときはクライアントのSWRが使われるはずなので、情報の更新間隔はクライアントのSWRのrevalidate間隔に従う
    revalidate: isStaticMode ? false : revalidateInterval,
  };
};

export const getStaticPaths: GetStaticPaths = async () => {
  const tasks = await fetchTasks().catch(() => []);

  return {
    paths: tasks.map((t) => ({ params: { id: t.id.toString() } })),
    fallback: isStaticMode ? false : "blocking",
  };
};

export default TasksDefault;
