package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var conn *websocket.Conn
var roc *exec.Cmd

func main() {
	cmd, err := exec.Command("rm", "-rf", "./downloadedMusic").CombinedOutput()
	cmdStr := strings.ToLower(string(cmd))
	fmt.Println(cmdStr)
	if err != nil {
		fmt.Println("Error deleting downloaded music: " + err.Error())
		return
	}
	cmd, err = exec.Command("rm", "-rf", "./pianoroll_rnn_nade").CombinedOutput()
	cmdStr = strings.ToLower(string(cmd))
	fmt.Println(cmdStr)
	if err != nil {
		fmt.Println("Error deleting previous neural network: " + err.Error())
		return
	}
	cmd, err = exec.Command("sudo", "apt", "install", "-y", "libasound2-dev").CombinedOutput()
	cmdStr = strings.ToLower(string(cmd))
	fmt.Println(cmdStr)
	if err != nil {
		fmt.Println("Error installing libasound2-dev: " + err.Error())
		return
	}
	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	r.Static("/magentarator", "./public")

	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	r.Run("localhost:9001")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func setupBeforeMagenta() {
	cmd, err := exec.Command("sudo", "apt", "install", "-y", "libasound2-dev").CombinedOutput()
	cmdStr := strings.ToLower(string(cmd))
	if err != nil {
		fmt.Println("Error installing libasound2-dev: " + err.Error())
		return
	}
	cmd, err = exec.Command("sudo", "apt", "install", "-y", "libjack-dev").CombinedOutput()
	cmdStr = strings.ToLower(string(cmd))
	if err != nil {
		fmt.Println("Error installing libjack-dev: " + err.Error())
		return
	}
	cmd, err = exec.Command("lshw", "-C", "display").CombinedOutput()
	cmdStr = strings.ToLower(string(cmd))
	if err != nil {
		fmt.Println("Error retrieving GPU info: " + err.Error())
		return
	}
	if strings.Contains(cmdStr, "nvidia") {
		fmt.Println("Detected NVIDIA GPU")
		fmt.Println("Checking for NVIDIA Drivers")
		cmd, _ := exec.Command("lsmod", "|", "grep", "nouveau").CombinedOutput()
		cmdStr := strings.ToLower(string(cmd))
		if strings.Contains(cmdStr, "nouveau") {
			fmt.Println("NVIDIA Drivers are not installed. Installing... (very risky)")
			_, err := exec.Command("yes", "|", "sudo", "ubuntu-drivers", "autoinstall").CombinedOutput()
			if err != nil {
				fmt.Println("Error installing NVIDIA Drivers: " + err.Error())
				return
			} else {
				setupBeforeMagenta()
				return
			}
		} else {
			fmt.Println("NVIDIA Drivers are likely enabled")
			cmd, err = exec.Command("nvcc", "-V").CombinedOutput()
			cmdStr := strings.ToLower(string(cmd))
			fmt.Println(cmdStr)
			if err == nil {
				fmt.Println("CUDA is installed")
				fmt.Println("Checking for Tensorflow")
				cmd, _ := exec.Command("python2", "-c", "'import tensorflow as tf; print(tf.__version__)'").CombinedOutput()
				cmdStr := strings.ToLower(string(cmd))
				if !strings.Contains(cmdStr, "Error") {
					fmt.Println("Tensorflow is already installed.")
					magentaSetup(true)
				} else {
					fmt.Println("Tensorflow is not installed")
					fmt.Println("Installing Tensorflow for GPU")
					cmd, _ := exec.Command("pip2", "install", "tensorflow-gpu").CombinedOutput()
					cmdStr := strings.ToLower(string(cmd))
					fmt.Println(cmdStr)
					if err != nil {
						fmt.Println("There was an error installing Tensorflow: " + err.Error())
						return
					} else {
						setupBeforeMagenta()
						return
					}
				}
			} else {
				fmt.Println("CUDA is not installed... Installing...")
				conn.WriteMessage(1, []byte("installingCuda"))
				cmd, err := exec.Command("sudo", "apt", "install", "nvidia-cuda-toolkit", "-y").CombinedOutput()
				cmdStr := strings.ToLower(string(cmd))
				fmt.Println(cmdStr)
				if err != nil {
					fmt.Println("Error installing CUDA: " + err.Error())
					return
				} else {
					fmt.Println("CUDA has been installed")
					setupBeforeMagenta()
					return
				}
				//cmdStr := strings.ToLower(string(cmd))
			}
		}
	} else {
		fmt.Println("No NVIDIA GPU Detected")
		fmt.Println("Checking if Tensorflow is installed")
		cmd, _ := exec.Command("python2", "-c", "'import tensorflow as tf; print(tf.__version__)'").CombinedOutput()
		cmdStr := strings.ToLower(string(cmd))
		if !strings.Contains(cmdStr, "Error") {
			fmt.Println("Tensorflow is already installed. Ending setup")
			magentaSetup(false)
		} else {
			fmt.Println("Tensorflow is not installed")
			fmt.Println("Installing Tensorflow for CPU")
			cmd, _ := exec.Command("pip2", "install", "tensorflow").CombinedOutput()
			cmdStr := strings.ToLower(string(cmd))
			fmt.Println(cmdStr)
			if err != nil {
				fmt.Println("There was an error installing Tensorflow: " + err.Error())
				return
			} else {
				fmt.Println("Tensorflow has been installed.")
				return
			}
		}
	}
}

func magentaSetup(gpu bool) {
	fmt.Println("Checking for Miniconda installation")
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir := usr.HomeDir
	fmt.Println(homeDir + "/.bashrc")
	cmd, err := exec.Command("source", homeDir+"/.bashrc").CombinedOutput()
	cmdStr := strings.ToLower(string(cmd))
	fmt.Println(cmdStr)
	cmd, err = exec.Command("conda", "--v").CombinedOutput()
	cmdStr = strings.ToLower(string(cmd))
	fmt.Println(cmdStr)
	if err != nil {
		fmt.Println("Miniconda is not installed, installing...")
		if strconv.IntSize == 64 {
			fmt.Println("Downloading Miniconda for 64-bit computers")
			err = downloadFile("conda.sh", "https://repo.continuum.io/miniconda/Miniconda2-latest-Linux-x86_64.sh")
			if err != nil {
				panic("Error downloading Miniconda: " + err.Error())
			} else {
				fmt.Println("Miniconda downloaded! Installing...")
				cmd, err := exec.Command("bash", "conda.sh", "-b", "-f", "-p", "$HOME/miniconda").CombinedOutput()
				cmdStr := strings.ToLower(string(cmd))
				if err != nil {
					fmt.Println("Error executing Miniconda setup file: " + err.Error())
				}
				fmt.Println(cmdStr)
				usr, err := user.Current()
				if err != nil {
					log.Fatal(err)
				}
				homeDir := usr.HomeDir
				fmt.Println(homeDir + "/.bashrc")
				f, err := os.OpenFile(homeDir+"/.bashrc", os.O_APPEND|os.O_WRONLY, 0644)
				_, err = f.WriteString("\nexport PATH=\"$HOME/miniconda/bin:$PATH\"")
				f.Close()
				if err != nil {
					fmt.Println("Error adding Miniconda to path: " + err.Error())
				}
				fmt.Println("Miniconda installed!")
				setupBeforeMagenta()
			}
		} else if strconv.IntSize == 32 {
			fmt.Println("Downloading Miniconda for 64-bit computers")
			err = downloadFile("conda.sh", "https://repo.continuum.io/miniconda/Miniconda2-latest-Linux-x86.sh")
			if err != nil {
				fmt.Println("Error downloading Miniconda: " + err.Error())
			} else {
				fmt.Println("Miniconda downloaded! Installing...")
				cmd, err := exec.Command("bash", "conda.sh", "-b", "-f", "-p", "$HOME/miniconda").CombinedOutput()
				cmdStr := strings.ToLower(string(cmd))
				if err != nil {
					fmt.Println("Error executing Miniconda setup file: " + err.Error())
				}
				fmt.Println(cmdStr)
				cmd, err = exec.Command("echo", "\"export PATH=\"$HOME/miniconda/bin:$PATH\"\"", ">>", "~/.bashrc").CombinedOutput()
				cmdStr = strings.ToLower(string(cmd))
				fmt.Println(cmdStr)
				if err != nil {
					fmt.Println("Error adding Miniconda to path: " + err.Error())
				}
				fmt.Println("Miniconda installed!")
				setupBeforeMagenta()
			}
		}
	}
	fmt.Println("Checking for a Magenta installation")
	cmd, err = exec.Command("source", "activate", "magenta").CombinedOutput()
	cmdStr = strings.ToLower(string(cmd))
	fmt.Println(cmdStr)
	if strings.Contains(cmdStr, "Error") {
		if gpu {
			fmt.Println("Magenta is not installed, installing Magenta with GPU support")
			cmd, err = exec.Command("conda", "remove", "-n", "magenta", "--all").CombinedOutput()
			if err != nil {
				fmt.Println("Error removing magenta environment")
			}
			cmd, err = exec.Command("conda", "create", "-y", "-f", "-n", "magenta", "python=2.7", "jupyter").CombinedOutput()
			cmdStr := strings.ToLower(string(cmd))
			fmt.Println(cmdStr)
			if err != nil {
				panic("There was an error creating the Magenta workspace: " + err.Error())
			} else {
				fmt.Println("Successfully created Magenta workspace, installing Magenta packages...")
				exec.Command("source", "activate", "magenta").CombinedOutput()
				cmd, err = exec.Command("pip", "install", "magenta-gpu").CombinedOutput()
				cmdStr := strings.ToLower(string(cmd))
				fmt.Println(cmdStr)
				if err != nil {
					panic("The magenta-gpu package could not be installed: " + err.Error())
				} else {
					fmt.Println("Magenta has been successfully installed. Ending setup")
					setupBeforeMagenta()
				}
			}
		} else {
			fmt.Println("Magenta is not installed, installing Magenta without GPU support")
			cmd, err = exec.Command("conda", "remove", "-n", "magenta", "--all").CombinedOutput()
			if err != nil {
				fmt.Println("Error removing magenta environment")
			}
			cmd, err = exec.Command("conda", "create", "-y", "-f", "-n", "magenta", "python=2.7", "jupyter").CombinedOutput()
			cmdStr := strings.ToLower(string(cmd))
			fmt.Println(cmdStr)
			if err != nil {
				panic("There was an error creating the Magenta workspace: " + err.Error())
			} else {
				fmt.Println("Successfully created Magenta workspace, installing Magenta packages...")
				exec.Command("source", "activate", "magenta").CombinedOutput()
				cmd, err = exec.Command("pip", "install", "magenta").CombinedOutput()
				cmdStr := strings.ToLower(string(cmd))
				fmt.Println(cmdStr)
				if err != nil {
					panic("The magenta package could not be installed: " + err.Error())
				} else {
					fmt.Println("Magenta has been successfully installed. Ending setup")
					setupBeforeMagenta()
				}
			}
		}
	} else {
		fmt.Println("Magenta is already installed, proceeding to main program")
		return
	}
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//type jsonNumber [][]

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, _ = wsupgrader.Upgrade(w, r, nil)
	setupBeforeMagenta()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		//fmt.Println(msg)
		if msg != nil {
			go func() {
				stringMsg := string(msg)
				//fmt.Println(stringMsg)
				if strings.HasPrefix(stringMsg, "o") {
					fmt.Println("Connection established")
					fmt.Println("Returning information about the computer")
					conn.WriteMessage(1, []byte("hey"))
				} else if strings.HasPrefix(stringMsg, "a") {
					stringMsg = strings.TrimLeft(stringMsg, "a")
					sliceMsg := strings.Split(stringMsg, ",")
					fmt.Println("Recieved artist(s): " + stringMsg)
					fmt.Println("Calling Python (beta)")
					for _, artist := range sliceMsg {
						cmd, err := exec.Command("python", "midiDownload.py", "a", artist).Output()
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(cmd)
					}
					dir, err := os.Getwd()
					roc = exec.Command("gnome-terminal", "--working-directory="+dir, "-e", "sh -c \"bash train.sh\"")
					cmd, err := roc.CombinedOutput()
					conn.WriteMessage(1, []byte("training"))
					cmdStr := strings.ToLower(string(cmd))
					fmt.Println(cmdStr)
					if err != nil {
						fmt.Println("Failed to train: " + err.Error())
					}
					go func() {
						for {
							json1 := new([][]float64)
							//getJSON("http://localhost:6006/data/plugin/scalars/scalars?run=run1%2Ftrain&tag=accuracy", json1)
							getJSON("http://localhost:6006/data/plugin/scalars/scalars?run=run1%2Ftrain&tag=metrics%2Faccuracy", json1)
							time.Sleep(time.Second)
							if len(*(json1)) > 0 {
								arr1 := (*json1)[len((*json1))-1]
								arr2 := arr1[len(arr1)-1]
								fmt.Println(string([]byte("p" + strconv.FormatFloat(arr2, 'f', -1, 64))))
								conn.WriteMessage(1, []byte("p"+strconv.FormatFloat(arr2, 'f', -1, 64)))
							}
							time.Sleep(time.Second)
						}
					}()
				} else if strings.HasPrefix(stringMsg, "g") {
					stringMsg = strings.TrimLeft(stringMsg, "g")
					sliceMsg := strings.Split(stringMsg, ",")
					fmt.Println("Recieved genre(s): " + stringMsg)
					fmt.Println("Calling Python (beta)")
					for _, genre := range sliceMsg {
						fmt.Println("Getting songs from " + genre)
						cmd, err := exec.Command("python2", "midiDownload.py", "g", genre).CombinedOutput()
						conn.WriteMessage(1, []byte("doneDownloading"))
						fmt.Println(string(cmd))
						if err != nil {
							fmt.Println(err)
							conn.WriteMessage(1, []byte("e"+err.Error()))
						} else {
							cmd, err := exec.Command("python2", "midiDownload.py", "g", genre).CombinedOutput()
							conn.WriteMessage(1, []byte("doneDownloading"))
							if err != nil {
								fmt.Println("Failed to download genres: " + err.Error())
							}
							fmt.Println(string(cmd))
						}
					}
					dir, err := os.Getwd()
					roc = exec.Command("gnome-terminal", "--working-directory="+dir, "-e", "sh -c \"bash train.sh\"")
					cmd, err := roc.CombinedOutput()
					conn.WriteMessage(1, []byte("training"))
					cmdStr := strings.ToLower(string(cmd))
					fmt.Println(cmdStr)
					if err != nil {
						fmt.Println("Failed to train: " + err.Error())
					}
					go func() {
						for {
							json1 := new([][]float64)
							//getJSON("http://localhost:6006/data/plugin/scalars/scalars?run=run1%2Ftrain&tag=accuracy", json1)
							getJSON("http://localhost:6006/data/plugin/scalars/scalars?run=run1%2Ftrain&tag=metrics%2Faccuracy", json1)
							time.Sleep(time.Second)
							if len(*(json1)) > 0 {
								arr1 := (*json1)[len((*json1))-1]
								arr2 := arr1[len(arr1)-1]
								fmt.Println(string([]byte("p" + strconv.FormatFloat(arr2, 'f', -1, 64))))
								conn.WriteMessage(1, []byte("p"+strconv.FormatFloat(arr2, 'f', -1, 64)))
							}
							time.Sleep(time.Second)
						}
					}()
				} else if strings.HasPrefix(stringMsg, "p") {
					stringMsg = strings.TrimLeft(stringMsg, "p")
					//sliceMsg := strings.Split(stringMsg, ",")
					fmt.Println(stringMsg)
				} else if strings.HasPrefix(stringMsg, "s") {
					stringMsg = strings.TrimLeft(stringMsg, "s")
					//sliceMsg := strings.Split(stringMsg, ",")
					fmt.Println(stringMsg)
				} else if strings.HasPrefix(stringMsg, "fgenerate") {
					roc.Process.Kill()
					cmd, err := exec.Command("kill", "-9", string(roc.Process.Pid)).CombinedOutput()
					if err != nil {
						fmt.Println("Couldn't kill it, using sudo")
						cmd, err = exec.Command("sudo", "kill", "-9", string(roc.Process.Pid)).CombinedOutput()
					}
					time.Sleep(time.Second * 30)
					dir, _ := os.Getwd()
					cmd, _ = exec.Command("gnome-terminal", "--working-directory="+dir, "-e", "sh -c \"bash generate.sh\"").CombinedOutput()
					conn.WriteMessage(1, []byte("generating"))
					cmdStr := strings.ToLower(string(cmd))
					fmt.Println(cmdStr)
				} else {
					fmt.Println("Couldn't understand: " + stringMsg)
				}
			}()
		}
		//conn.WriteMessage(t, msg)
	}
}
