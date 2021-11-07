import {
  Modal,
  ModalBody,
  ModalContent,
  ModalOverlay,
  useDisclosure,
} from "@chakra-ui/react";
import TaskModalBody from "components/taskmodalbody";
import { useRouter } from "next/router";
import Loading from "../../components/loading";
import Tasks from "../../components/tasks";
import useAccount from "../../lib/api/account";
import useTasks from "../../lib/api/tasks";
import parentpath from "../../lib/parentpath";

const TasksDefault = () => {
  const router = useRouter();
  const { id } = router.query;
  const { onClose } = useDisclosure();

  const { data: tasks } = useTasks();
  const { data: account } = useAccount();

  if (!tasks || account === undefined) {
    return <Loading />;
  }

  const filterdTasks = tasks.filter((t) => t.id === Number(id));
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
        size="xl"
      >
        <ModalOverlay />
        <ModalContent>
          <ModalBody>
            <TaskModalBody task={task} />
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
};

export default TasksDefault;
