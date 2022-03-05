import {
  Modal,
  ModalBody,
  ModalContent,
  ModalOverlay,
  useDisclosure,
} from "@chakra-ui/react";
import TaskModalBody from "components/taskmodalbody";
import { isStaticMode } from "lib/static";
import { GetStaticPaths, GetStaticProps } from "next";
import { useRouter } from "next/router";
import Loading from "../../components/loading";
import Tasks from "../../components/tasks";
import useAccount, { Account, fetchAccount } from "../../lib/api/account";
import useTasks, { fetchTasks, Task } from "../../lib/api/tasks";
import parentpath from "../../lib/parentpath";

interface TaskProps {
  taskID: number;
  tasks: Task[];
  account: Account | null;
}

const TasksDefault = ({
  taskID,
  tasks: defaultTasks,
  account: defaultAccount,
}: TaskProps) => {
  const router = useRouter();
  const { onClose } = useDisclosure();

  const { data: tasks } = useTasks(defaultTasks);
  const { data: account } = useAccount(defaultAccount);

  if (!tasks || account === undefined) {
    return <Loading />;
  }

  const filterdTasks = tasks.filter((t) => t.id === taskID);
  if (filterdTasks.length !== 1) {
    return <Loading />;
  }
  const task = filterdTasks[0];

  return (
    <>
      <Tasks tasks={tasks} team={account} />
      <Modal
        isOpen={true}
        onClose={() => {
          onClose();
          router.push(parentpath(router.route), undefined, {
            scroll: false,
            shallow: true,
          });
        }}
        size="4xl"
      >
        <ModalOverlay />
        <ModalContent
          sx={{
            backgroundColor: "#ffffff",
          }}
        >
          <ModalBody>
            <TaskModalBody task={task} />
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
};

export const getStaticProps: GetStaticProps<TaskProps> = async (context) => {
  const id = context.params?.id;
  const account = isStaticMode ? null : await fetchAccount();
  const tasks = await fetchTasks();
  return {
    props: {
      taskID: Number(id),
      tasks: tasks,
      account: account,
    },
  };
};

export const getStaticPaths: GetStaticPaths = async () => {
  const tasks = await fetchTasks();

  return {
    paths: tasks.map((t) => ({ params: { id: t.id.toString() } })),
    fallback: false,
  };
};

export default TasksDefault;
