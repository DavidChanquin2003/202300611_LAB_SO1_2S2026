#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("David Arturo Chanquin Velasquez");
MODULE_DESCRIPTION("Modulo de telemetria de RAM y procesos para SO1");

#define PROCFS_NAME "continfo_pr2_so1_202300611"

static int escribir_archivo(struct seq_file *archivo, void *v) {
    struct sysinfo i;
    si_meminfo(&i);

    // Vamos calculando en Megabytes 
    long total_ram = (i.totalram * i.mem_unit) / (1024 * 1024);
    long free_ram = (i.freeram * i.mem_unit) / (1024 * 1024);
    long used_ram = total_ram - free_ram;

    seq_printf(archivo, "Total de memoria RAM: %ld MB\n", total_ram);
    seq_printf(archivo, "Memoria RAM libre: %ld MB\n", free_ram);
    seq_printf(archivo, "Memoria RAM en uso: %ld MB\n", used_ram);
    seq_printf(archivo, "-----------------------------------------\n");
    
    // NOTA: Aqui insertaremos la iteracion de task_struct despues
    
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
    printk(KERN_INFO "Modulo de RAM cargado en /proc/%s\n", PROCFS_NAME);
    return 0;
}

static void __exit modulo_cleanup(void) {
    remove_proc_entry(PROCFS_NAME, NULL);
    printk(KERN_INFO "Modulo descargado\n");
}

module_init(modulo_init);
module_exit(modulo_cleanup);