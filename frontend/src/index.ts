import YQ from './yq/yq';
import mdui from 'mdui';

class Index {
    oss: string = 'https://tongdy-mediaen.oss-cn-beijing.aliyuncs.com/'
    api: string = 'https://tool.mytongdy.com/api/ossdirview/tongdy_dirviewer/';
    tableLineCode: string = '';
    fileListTbody: HTMLElement = document.getElementById('fileListTbody')!
    t1: HTMLElement = document.getElementById('t1')!
    t2: HTMLElement = document.getElementById('t2')!
    btnBack: HTMLElement = document.getElementById('btnBack')!
    btnReload: HTMLElement = document.getElementById('btnReload')!
    title: string = document.title;

    constructor() {
        YQ.get('filelist.template.html', undefined, (data: XMLHttpRequest | null, status: number) => {
            if (data) {
                this.tableLineCode = data.responseText;
            }
        }, false);
        this.btnBack.addEventListener('click', () => {
            history.go(-1);
        });
        this.btnReload.addEventListener('click', () => {
            location.reload();
        });
        this.reloadData();
    }

    reloadData() {
        const id: string = YQ.argv('id') as string;
        const secret: string = YQ.argv('secret') as string;
        const path: string = YQ.argv('path') as string;
        this.t2.innerText = '/' + id + '/' + path;
        document.title = this.title + ' ' + this.t2.innerText;
        if (id.length == 0 || secret.length == 0) {
            this.fileListTbody.innerHTML = '<p><center>none</center></p>';
            return;
        }
        this.fileListTbody.innerHTML = '<p><center><div class="mdui-spinner mdui-spinner-colorful"></div></center></p>';
        mdui.mutation();
        YQ.get(this.api, {
            id: id,
            secret: secret,
            path: path,
        }, (data: XMLHttpRequest | null, status: number) => {
            if (!data) {
                return;
            }
            const jsonData: any[] = JSON.parse(data.responseText).data as any[];
            let html = '';
            for (const key in jsonData) {
                if (Object.prototype.hasOwnProperty.call(jsonData, key)) {
                    const file: any = jsonData[key];
                    const name: string = file.path;
                    const size: number = file.size;
                    const type: string = file.type;
                    const extname: string = file.path.split('.').pop();
                    const isDir: boolean = (file.type == 'folder');
                    const icon: string = isDir ? 'folder' : this.fileicon(extname);
                    const rowurl: string = isDir ? window.location.pathname + '#id=' + id + '&secret=' + secret + '&path=' + file.path : this.oss + id + '/' + path + file.path;
                    html += YQ.loadTemplateHtml(this.tableLineCode, undefined, [
                        ['icons', icon],
                        ['filepath', rowurl],
                        ['atarget', isDir ? '' : '_blank'],
                        ['filename', name],
                        ['filetype', (isDir ? '' : extname + ' ') + type],
                        ['filesize', isDir ? '' : this.sizeConver(size)],
                    ]);
                }
            }
            this.fileListTbody.innerHTML = html;
            mdui.updateTables();
        }, false);
    }

    fileicon(extname: String) {
        if (extname == 'jpg' || extname == 'png' || extname == 'svg' || extname == 'tif' || extname == 'jiff' || extname == 'tiff' || extname == 'bmp') {
            return 'image';
        } else if (extname == 'zip' || extname == '7z' || extname == 'tar' || extname == 'xz' || extname == 'rar' || extname == 'cab') {
            return 'archive';
        } else {
            return 'insert_drive_file';
        }
    }

    sizeConver(limit: number) {
        let size: string = '';
        if (limit < 0.1 * 1024) {
            size = limit.toFixed(2) + ' B';
        } else if (limit < 0.1 * 1024 * 1024) {
            size = (limit / 1024).toFixed(2) + ' KB';
        } else if (limit < 0.1 * 1024 * 1024 * 1024) {
            size = (limit / (1024 * 1024)).toFixed(2) + ' MB';
        } else {
            size = (limit / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
        }
        const sizestr: string = size;
        const len: number = sizestr.indexOf('\.');
        const dec: string = sizestr.substr(len + 1, 2);
        if (dec == '00') {
            return sizestr.substring(0, len) + sizestr.substr(len + 3, 2);
        }
        return sizestr;
    }
}

window.onload = () => {
    new Index();
};

window.onhashchange = () => {
    location.reload();
};
