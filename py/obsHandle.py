import os
import logging
from obs import ObsClient

logging.basicConfig(level=logging.NOTSET)


class OBSHandler:
    def __init__(self, obs_ak, obs_sk, obs_bucketname, obs_endpoint):
        self.access_key = obs_ak
        self.secret_key = obs_sk
        self.bucket_name = obs_bucketname
        self.endpoint = obs_endpoint
        self.server = obs_endpoint
        self.obs_client = self.init_obs()
        self.maxkeys = 5000  # 查询的对象最大个数

    # 初始化obs
    def init_obs(self):
        obs_client = ObsClient(
            access_key_id=self.access_key,
            secret_access_key=self.secret_key,
            server=self.server
        )
        return obs_client

    def close_obs(self):
        self.obs_client.close()

    def read_file(self, path):
        """
        二进制读取配置文件
        :param path:
        :return:
        """
        try:
            resp = self.obs_client.getObject(
                self.bucket_name, path, loadStreamInMemory=True)
            if resp.status < 300:
                # 获取对象内容
                return {
                    "status": 200,
                    "msg": "获取配置文件成功",
                    "content": bytes.decode(resp.body.buffer, "utf-8"),
                    "size": resp.body.size
                }
            else:
                return {
                    "status": -1,
                    "msg": "获取失败，失败码: %s\t 失败消息: %s" % (resp.errorCode, resp.errorMessage),
                    "content": "",
                    "size": 0
                }
        except:
            logging.error('obs_client读取文件失败')


if __name__ == "__main__":
    obs = OBSHandler()
    # 关闭obs_client
    obs.close_obs()
