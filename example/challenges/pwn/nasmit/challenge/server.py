#!/usr/bin/env python3
import os
import string
import random
import subprocess

def check_syntax(asm):
    banned = {
        '%': 'Macro is not allowed',
        'syscall': 'System call is not allowed',
        'int'    : 'Interruption is not allowed',
        'incbin' : 'External file is not allowed',
        'extern' : 'External functon is not allowed',
        ':': 'Label is not allowed',
        '$': 'Special token is not allowed'
    }

    for inst in asm:
        for keyword in banned:
            if keyword in inst:
                return banned[keyword]

    return None

def assemble(asm):
    def randstr(l):
        return ''.join([random.choice(string.ascii_letters) for i in range(l)])

    path_asm = '/tmp/' + randstr(16) + '.S'
    path_obj = '/tmp/' + randstr(16) + '.o'
    path_elf = '/tmp/' + randstr(16) + '.elf'
    with open(path_asm, 'w') as f:
        f.writelines([
            'default rel\n'
            '%macro exit 1\n',
            '  mov rdi, %1\n',
            '  jmp _exit WRT ..plt\n',
            '%endmacro\n',
            'extern _exit\n'
            'section .text\n',
            'global _start\n',
            '_start:\n'
        ])
        f.write('\n'.join(asm))

    # assemble
    p = subprocess.Popen(['timeout', '1',
                          'nasm', path_asm, '-o', path_obj, '-fELF64'],
                         stderr=subprocess.PIPE)
    err = p.stderr.read().decode()
    os.remove(path_asm)
    if err:
        return path_obj, err

    # link
    p = subprocess.Popen(['timeout', '1',
                          'ld',
                          '-dynamic-linker', '/lib64/ld-linux-x86-64.so.2',
                          '-lc', path_obj, '-o', path_elf,
                          '-znoexecstack', '-znow', '-zrelro', '-pie'],
                         stderr=subprocess.PIPE)
    err = p.stderr.read().decode()
    os.remove(path_obj)
    if err:
        return path_elf, err
    
    return path_elf, None

def run(path):
    p = subprocess.Popen(['timeout', '-sSIGKILL', '1', path],
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    while True:
        ret = p.poll()
        if ret is not None:
            break

    return ret

def main():
    print("Enter your assembly (Empty line to end)")

    # Read input
    asm = []
    for i in range(100):
        line = input()
        if line == '':
            break
        else:
            asm.append(line)

    # Check security syntax
    err = check_syntax(asm)
    if err:
        print("[ERROR] Security check:")
        print(err)
        return

    # Assemble and link
    path, err = assemble(asm)
    if err:
        print("[ERROR] Assembler error:")
        print(err)
        os.remove(path)
        return

    # Run
    ret = run(path)
    os.remove(path)

    print("[INFO] Status code:")
    print(ret)

if __name__ == '__main__':
    try:
        main()
    except Exception as e:
        print("[ERROR] Unhandled exception:")
        print(e)
