import {
  Modal,
  ModalBody,
  ModalContent,
  ModalOverlay,
  useDisclosure,
} from "@chakra-ui/react";
import TaskModalBody from "components/taskmodalbody";
import { useRouter } from "next/router";
import Loading from "../../../components/loading";
import parentpath from "../../../lib/parentpath";
import Tasks from "./index";
import useAdminTasks from "../../../lib/api/admin/tasks";

const TasksDefault = () => {
  const router = useRouter();
  const { id } = router.query;
  const { onClose } = useDisclosure();

  const { data: tasks } = useAdminTasks();
  if (!tasks) {
    return <Loading />;
  }

  const filterdTasks = tasks.filter((t) => t.id === Number(id));
  if (filterdTasks.length !== 1) {
    return <Loading />;
  }
  const task = filterdTasks[0];

  return (
    <>
      <Tasks />
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

export default TasksDefault;
