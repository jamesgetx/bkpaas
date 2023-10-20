# Generated by Django 3.2.12 on 2023-10-19 03:07

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    dependencies = [
        ('bkapp_model', '0003_auto_20231017_1501'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='moduleprocessspec',
            name='target_status',
        ),
        migrations.AlterField(
            model_name='moduleprocessspec',
            name='plan_name',
            field=models.CharField(help_text='仅存储方案名称', max_length=32),
        ),
        migrations.AlterField(
            model_name='moduleprocessspec',
            name='scaling_config',
            field=models.JSONField(null=True, verbose_name='自动扩缩容配置'),
        ),
        migrations.CreateModel(
            name='ProcessSpecEnvOverlay',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('region', models.CharField(help_text='部署区域', max_length=32)),
                ('created', models.DateTimeField(auto_now_add=True)),
                ('updated', models.DateTimeField(auto_now=True)),
                ('environment_name', models.CharField(choices=[('stag', '预发布环境'), ('prod', '生产环境')], max_length=16, verbose_name='环境名称')),
                ('target_replicas', models.IntegerField(default=1, verbose_name='期望副本数')),
                ('plan_name', models.CharField(help_text='仅存储方案名称', max_length=32)),
                ('autoscaling', models.BooleanField(default=False, verbose_name='是否启用自动扩缩容')),
                ('scaling_config', models.JSONField(null=True, verbose_name='自动扩缩容配置')),
                ('proc_spec', models.ForeignKey(db_constraint=False, on_delete=django.db.models.deletion.CASCADE, related_name='env_overlay', to='bkapp_model.moduleprocessspec')),
            ],
            options={
                'unique_together': {('proc_spec', 'environment_name')},
            },
        ),
    ]