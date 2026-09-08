#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched/signal.h>
#include <linux/sched.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("David Arturo Chanquin Velasquez");
MODULE_DESCRIPTION("Modulo de telemetria de RAM y procesos para SO1");

#define PROCFS_NAME "continfo_pr2_so1_202300611"

static int escribir_archivo(struct seq_file *archivo, void *v) {
    struct sysinfo i;
    struct task_struct *task;
    unsigned long vsz, rss, mem_porcentaje;
    long total_ram_kb;

    si_meminfo(&i);

    // Calculando en Megabytes para RAM General
    long total_ram_mb = (i.totalram * i.mem_unit) / (1024 * 1024);
    long free_ram_mb = (i.freeram * i.mem_unit) / (1024 * 1024);
    long used_ram_mb = total_ram_mb - free_ram_mb;
    
    total_ram_kb = (i.totalram * i.mem_unit) / 1024;

    seq_printf(archivo, "Total de memoria RAM: %ld MB\n", total_ram_mb);
    seq_printf(archivo, "Memoria RAM libre: %ld MB\n", free_ram_mb);
    seq_printf(archivo, "Memoria RAM en uso: %ld MB\n", used_ram_mb);
    seq_printf(archivo, "-----------------------------------------\n");
    
    seq_printf(archivo, "PID | Nombre | VSZ (KB) | RSS (KB) | %% RAM | CPU (Jiffies)\n");

    // Iterar sobre todos los procesos del sistema
    for_each_process(task) {
        if (task->mm) {
            vsz = task->mm->total_vm << (PAGE_SHIFT - 10);
            rss = get_mm_rss(task->mm) << (PAGE_SHIFT - 10);
            
            // Calculo de porcentaje de memoria
            mem_porcentaje = (total_ram_kb > 0) ? (rss * 100) / total_ram_kb : 0;
        } else {
            vsz = 0;
            rss = 0;
            mem_porcentaje = 0;
        }

        seq_printf(archivo, "%d | %s | %lu | %lu | %lu%% | %llu\n",
                   task->pid, 
                   task->comm, 
                   vsz, 
                   rss, 
                   mem_porcentaje, 
                   task->utime); // utime significa o representa el tiempo en espacio de usuario
    }
    
    return 0;
}

static int al_abrir(struct inode *inode, struct file *file) {
    return single_open(file, escribir_archivo, NULL);
}

static const struct proc_ops operaciones = {
    .proc_open = al_abrir,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

static int __init modulo_init(void) {
    proc_create(PROCFS_NAME, 0, NULL, &operaciones);
    printk(KERN_INFO "Modulo cargado exitosamente.\n");
    return 0;
}

static void __exit modulo_cleanup(void) {
    remove_proc_entry(PROCFS_NAME, NULL);
    printk(KERN_INFO "Modulo descargado exitosamente.\n");
}

module_init(modulo_init);
module_exit(modulo_cleanup);