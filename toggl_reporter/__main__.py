import sys


def main():
    try:
        from .cli import run  # noqa: WPS300, WPS433
        exit_status = run()
    except KeyboardInterrupt:
        # 128+2 SIGINT
        # <http://www.tldp.org/LDP/abs/html/exitcodes.html>
        exit_status = 130

    sys.exit(exit_status)


if __name__ == '__main__':
    main()
