import { fetchCTF } from "lib/api/ctf";
import { AllPageProps } from "lib/pages";
import { isStaticMode } from "lib/static";
import { GetStaticPaths, GetStaticProps } from "next";
import { useRouter } from "next/router";
import Loading from "../../components/loading";
import useAccount, { Account, fetchAccount } from "../../lib/api/account";
import useTasks, { fetchTasks, Task } from "../../lib/api/tasks";
import parentpath from "../../lib/parentpath";
import TaskView from "theme/task";
import useMessage from "lib/useMessage";
import { FlagSubmitParams } from "props/task";
import { useForm, SubmitHandler } from "react-hook-form";
import { api } from "lib/api";
import { filterTask, sortTask, isSolved } from "lib/tasks";
import { useCallback } from "react";

type taskProps = {
  taskID: number;
  tasks: Task[];
  account: Account | null;
} & AllPageProps;

const TasksDefault = ({
  taskID,
  tasks: defaultTasks,
  account: defaultAccount,
}: taskProps) => {
  const router = useRouter();
  const { data: tasks, mutate } = useTasks(defaultTasks);
  const { data: account } = useAccount(defaultAccount);
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

  return (
    <>
      <TaskView
        task={task}
        tasks={tasks}
        tasksPath={parentpath(router.pathname)}
        registerFlag={register}
        onFlagSubmit={handleSubmit(onSubmit)}
        filterTask={filterTask}
        sortTask={sortTask}
        isSolved={isSolvedByTeam}
      />
    </>
  );
};

export const getStaticProps: GetStaticProps<taskProps> = async (context) => {
  const id = context.params?.id;
  const account = isStaticMode ? null : await fetchAccount();
  const tasks = await fetchTasks();
  const ctf = await fetchCTF();
  return {
    props: {
      taskID: Number(id),
      tasks: tasks,
      account: account,
      ctf: ctf,
    },
    revalidate: isStaticMode ? undefined : 30, // revalidate every 1 seconds
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
