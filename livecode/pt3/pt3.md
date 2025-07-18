## Сериализация
Виды сериализации:
1. CSV/TSV
2. JSON
3. msgp 
4. arrow


- ls /  абсолютные пути (самый корень файловой системы)
- ls .  относительные пути (начинаются с точки .  типо точка текущая директория) показывает то, где я щас нахожусь
Когда мы что-то записываем в файл мы должны указать путь к нему: либо абсолютный, либо относительный


go build .  с помощью этой команды создается бинарный файл программы 
этот файл можно запустить
например,
```
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo$ ./key-value_project/myproj
/mnt/c/ProjectsGOLANG/projgo
[{2  started 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC {before 28.10 now}} {3  not started 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC { }}]
```
относительные пути нужно использовать очень осторожно 

чтобы не решать проблему с путями 

var rootDir string 
filepath.Join("/home/student/", "tasks.json")


err = os.WriteFile("tasks.json", b, 755) // 755 - разрешение
это ненадежно. Во-первых, пока вы пишите эти байты: у вас могут быть:
- огромные данные (гигабайты)
- на середине сломаться процесс
и получится то, что в файле останется половина того, что вы написали 
может быть такое, что кто-то другой начал писать этот файл
Суть в том, что запись в файлы это не атомартная операция 

В линкусе есть подход, как атомарно написать файл
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ ls
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ echo "data"
data
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ echo "data" >
bash: syntax error near unexpected token `newline'
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ echo "data" > file.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ cat file.txt
data
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ echo "duck" > file.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ cat
cat           catchsegv     catman        catsrv.dll    catsrvps.dll  catsrvut.dll  
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ cat
cat           catchsegv     catman        catsrv.dll    catsrvps.dll  catsrvut.dll  
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ cat file.txt
duck
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ echo "duck" >> file.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ cat file.txt
duck
duck
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ ls
file.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ cp file.txt file2.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ ls
file.txt  file2.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ cat file2.txt
duck
duck
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ rm file2.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ mv file.txt file2.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ ls
file2.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ rm ./file2.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ echo "duck" > file.txt
student@DESKTOP-QNI8A2Q:/mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test$ mv file.txt file2.txt